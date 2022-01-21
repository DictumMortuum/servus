package food

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DictumMortuum/servus/pkg/config"
	"github.com/DictumMortuum/servus/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/microcosm-cc/bluemonday"
	"io/ioutil"
	"net/http"
	"strings"
)

func Scrape(c *gin.Context) {
	url := c.PostForm("url")

	data, err := scrapeAkis(url)
	if err != nil {
		util.Error(c, err)
		return
	}

	util.Success(c, &data)
	return
}

type AkisJSON struct {
	Props struct {
		PageProps struct {
			Recipe struct {
				Data struct {
					Nutrition struct {
						Per      string                   `json:"nutrition_per"`
						Sections []map[string]interface{} `json:"sections"`
					} `json:"nutrition"`
					Slug   string `json:"slug"`
					Video  string `json:"video_url"`
					Method []struct {
						Section string              `json:"section"`
						Steps   []map[string]string `json:"steps"`
					} `json:"method"`
					IngredientSections []struct {
						Ingredients []struct {
							Quantity string `json:"quantity"`
							Title    string `json:"title"`
							Unit     string `json:"unit"`
						} `json:"ingredients"`
					} `json:"ingredient_sections"`
				} `json:"data"`
			} `json:"ssRecipe"`
		} `json:"pageProps"`
	} `json:"props"`
}

func scrapeAkis(url string) (map[string]interface{}, error) {
	p := bluemonday.NewPolicy()

	rs := map[string]interface{}{}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"),
	)
	c.WithTransport(t)

	rs["tags"] = []string{}
	rs["recipe_category"] = []string{}
	rs["notes"] = []string{}

	c.OnHTML(".assets-wrapper", func(e *colly.HTMLElement) {
		rs["name"] = e.ChildText(".assets-wrapper h1")
		rs["recipe_yield"] = e.ChildText(".spec-item.portions")
	})

	c.OnHTML("meta", func(e *colly.HTMLElement) {
		key := e.Attr("name")

		if key == "name" || key == "description" {
			val := e.Attr("content")
			rs[key] = val
		}

		prop := e.Attr("property")
		if prop == "og:image" {
			rs["image"] = e.Attr("content")
		}
	})

	c.OnHTML(".time-specs", func(e *colly.HTMLElement) {
		rs["prepTime"] = e.ChildText(".preparation .content")
		rs["totalTime"] = e.ChildText(".rest .content")
		rs["performTime"] = e.ChildText(".cooking .content")
	})

	c.OnHTML("#__NEXT_DATA__", func(e *colly.HTMLElement) {
		var tmp AkisJSON
		json.Unmarshal([]byte(e.Text), &tmp)
		rs["slug"] = "akis-" + tmp.Props.PageProps.Recipe.Data.Slug
		rs["org_url"] = tmp.Props.PageProps.Recipe.Data.Video
		rs["json"] = tmp

		instructions := []map[string]string{}

		for _, meth := range tmp.Props.PageProps.Recipe.Data.Method {
			for _, item := range meth.Steps {
				instruction := map[string]string{}
				instruction["text"] = p.Sanitize(item["step"])
				instructions = append(instructions, instruction)
			}
		}

		ingredients := []string{}

		for _, section := range tmp.Props.PageProps.Recipe.Data.IngredientSections {
			for _, item := range section.Ingredients {
				ingredients = append(ingredients, strings.TrimSpace(item.Quantity+" "+item.Unit+" "+item.Title))
			}
		}

		// "nutrition": {
		// 	"calories": null,
		// 	"fatContent": null,
		// 	"proteinContent": null,
		// 	"carbohydrateContent": null,
		// 	"fiberContent": null,
		// 	"sodiumContent": null,
		// 	"sugarContent": null
		// },

		rs["recipe_instructions"] = instructions
		rs["recipe_ingredient"] = ingredients
	})

	c.Visit(url)

	slug, err := createRecipe(rs)
	if err != nil {
		return nil, err
	}

	_, err = createRecipeImage(slug, rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func createRecipe(payload map[string]interface{}) (*string, error) {
	jsonStr, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://mealie.dictummortuum.com/api/recipes/create", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.App.Mealie.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tmp := string(body)

	return &tmp, nil
}

func createRecipeImage(slug *string, payload map[string]interface{}) (map[string]interface{}, error) {
	jsonStr, _ := json.Marshal(map[string]interface{}{
		"url": payload["image"],
	})

	fmt.Println(string(jsonStr), fmt.Sprintf("https://mealie.dictummortuum.com/api/recipes/%s/image", strings.Trim(*slug, "\"")))

	req, err := http.NewRequest("POST", fmt.Sprintf("https://mealie.dictummortuum.com/api/recipes/%s/image", strings.Trim(*slug, "\"")), bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.App.Mealie.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
