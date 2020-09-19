package main

import (
	"fmt"
	"github.com/IBM/go-sdk-core/core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jeffotoni/gconcat"
	"github.com/watson-developer-cloud/go-sdk/naturallanguageunderstandingv1"
	"github.com/watson-developer-cloud/go-sdk/speechtotextv1"
	"io"


	"math"
	"os"
	"strings"
)

type Dados struct {
	Car  string `form:"car"`
	Text string `form:"text"`
}

func main() {
	var urlnnul string
	app := fiber.New(fiber.Config{BodyLimit: 12 * 1024 * 1024})
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${ip} ${path} ${time} \n",
		TimeFormat: "02-Jan-2006",
		Output:     os.Stdout}))

	app.Post("/behinthecode8", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		println("INICIANDO")
		d := new(Dados)
		err := c.BodyParser(d)
		switch err {
		case nil:
			fmt.Println(d)
			break
		default:
			return c.Status(200).SendString(`{"recommendation": "","entities": []}`)
		}
		switch d.Text {
		case "":
			println("INICIANDO SALVAMENTO DO ARQUIVO DE AUDIO")
			file, err := c.FormFile("audio")
			switch err {
			case nil:
				err := c.SaveFile(file, gconcat.Build("./audio/", file.Filename))
				switch err {
				case nil:
					break
				default:
					fmt.Println("ERROR:", err)
					return c.Status(200).SendString(`{"recommendation": "","entities": []}`)
				}
			default:
				fmt.Println("ERROR:", err)
				return c.Status(200).SendString(`{"recommendation": "","entities": []}`)
			}
			println("ENVIANDO PARA O SST")
			authenticator := &core.IamAuthenticator{
				ApiKey: os.Getenv("APIKEY_BEHINDCODE8"),
			}

			options := &speechtotextv1.SpeechToTextV1Options{
				Authenticator: authenticator,
			}

			speechToText, speechToTextErr := speechtotextv1.NewSpeechToTextV1(options)

			if speechToTextErr != nil {
				return c.Status(500).SendString("deu ruim")
			}

			speechToText.SetServiceURL(os.Getenv("URL1_BEHINDCODE8"))
			file1,err:=file.Open()
			if err!=nil{
				return c.Status(500).SendString("deu ruim")
			}
			var audioFile io.ReadCloser
			audioFile=file1

			result, _, responseErr := speechToText.Recognize(
				&speechtotextv1.RecognizeOptions{
					Audio:       audioFile,
					ContentType: core.StringPtr("application/octet-stream"),
					Model:       core.StringPtr("pt-BR_BroadbandModel"),
				},
			)
			if responseErr != nil {
				return c.Status(500).SendString(responseErr.Error())
			}

			for i, result1 := range result.Results {
				urlnnul = gconcat.Build(urlnnul, *result1.Alternatives[i].Transcript, "\n")
				fmt.Println(*result1.Alternatives[i].Transcript, "\n")
			}
			break
		default:
			urlnnul = d.Text
		}

		println("ENVIANDO PARA O NLU")
		authenticator1 := &core.IamAuthenticator{
			ApiKey: os.Getenv("APIKEY1_BEHINDCODE8"),
		}

		options1 := &naturallanguageunderstandingv1.NaturalLanguageUnderstandingV1Options{
			Version:       "2020-09-17",
			Authenticator: authenticator1,
		}

		naturalLanguageUnderstanding, naturalLanguageUnderstandingErr := naturallanguageunderstandingv1.NewNaturalLanguageUnderstandingV1(options1)

		if naturalLanguageUnderstandingErr != nil {
			fmt.Println("ERROR:", naturalLanguageUnderstandingErr)
			return c.Status(500).SendString("deu ruim")
		}
		naturalLanguageUnderstanding.Service.SetServiceURL(os.Getenv("URL2_BEHINDCODE8"))
		id := os.Getenv("MODELOID")

		result1, _, responseErr1 := naturalLanguageUnderstanding.Analyze(
			&naturallanguageunderstandingv1.AnalyzeOptions{
				Text: &urlnnul,
				Features: &naturallanguageunderstandingv1.Features{
					Entities: &naturallanguageunderstandingv1.EntitiesOptions{
						Mentions:  core.BoolPtr(true),
						Model:     core.StringPtr(id),
						Sentiment: core.BoolPtr(true),
					},
				},
			},
		)
		switch responseErr1 {
		case nil:
			tamanho := len(result1.Entities)
			switch tamanho {
			case 0:
				return c.Status(200).SendString(`{"recommendation": "","entities": []}`)
			}
			var soma float64
			var maior float64
			var maior2 float64
			var entidade string
			var entidade2 string
			recomendacao := `{"recommendation":`
			retorno := `"entities": [`
			for i :=range result1.Entities {
				retorno = gconcat.Build(retorno,`{"entity": "`, *result1.Entities[i].Type, `","sentiment":`, *result1.Entities[i].Sentiment.Score, `,"mention": "`,*result1.Entities[i].Text,`"},`)
				soma = soma + *result1.Entities[i].Sentiment.Score
				switch {
				case i==0:
					maior = *result1.Entities[i].Sentiment.Score
					entidade = *result1.Entities[i].Type
				case *result1.Entities[i].Sentiment.Score > maior :
					print("to aqui 1")
					maior2 = maior
					maior = *result1.Entities[i].Sentiment.Score
					entidade2 = entidade
					entidade = *result1.Entities[i].Type
					break
				}
			}
			retorno = strings.Trim(retorno, ",")
			retorno = gconcat.Build(retorno, `]}`)
			type Recomen struct {
				r1 string
				r2 string
			}

			var m map[string]Recomen
			m = make(map[string]Recomen)
			m["SEGURANCA"] = Recomen{
				"TOURO", "ARGO",
			}
			m["CONSUMO"] = Recomen{
				"FIORINO", "FIAT 500",
			}
			m["DESEMPENHO"] = Recomen{
				"MAREA", "RENEGADE",
			}
			m["MANUTENCAO"] = Recomen{
				"FIORINO ", "LINEA",
			}
			m["CONFORTO"] = Recomen{
				"RENEGADE", "TORO",
			}
			m["DESIGN"] = Recomen{
				"TORO", "DUCATO",
			}
			m["ACESSORIOS"] = Recomen{
				"CRONOS", "TORO",
			}

			switch {
			case maior-maior2 == 0:
				switch {
				case entidade == "SEGURANCA":
					switch d.Car {
					case m["SEGURANCA"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r2, `",`)
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r1, `",`)
					}
					break
				case entidade2 == "SEGURANCA":
					switch d.Car {
					case m["SEGURANCA"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r2, `",`)
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r1, `",`)
					}
					break
				case entidade == "CONSUMO":
					switch d.Car {
					case m["CONSUMO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r2, `",`)
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r1, `",`)
					}
					break
				case entidade2 == "CONSUMO":
					switch d.Car {
					case m["CONSUMO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r1, `",`)
						break
					}
					break
				case entidade == "DESEMPENHO":
					switch d.Car {
					case m["DESEMPENHO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r1, `",`)
						break
					}
					break
				case entidade2 == "DESEMPENHO":
					switch d.Car {
					case m["DESEMPENHO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r1, `",`)
						break
					}
					break
				case entidade == "MANUTENCAO":
					switch d.Car {
					case m["MANUTENCAO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r1, `",`)
						break
					}
					break
				case entidade2 == "MANUTENCAO":
					switch d.Car {
					case m["MANUTENCAO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r1, `",`)
						break
					}
					break
				case entidade == "CONFORTO":
					switch d.Car {
					case m["CONFORTO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r1, `",`)
						break
					}
					break
				case entidade2 == "CONFORTO":
					switch d.Car {
					case m["CONFORTO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r1, `",`)
						break
					}
					break
				case entidade == "DESIGN":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r1, `",`)
						break
					}
					break
				case entidade2 == "DESIGN":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r1, `",`)
						break
					}
					break
				case entidade == "ACESSORIOS":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r1, `",`)
						break
					}
					break
				case entidade2 == "ACESSORIOS":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r1, `",`)
						break
					}
					break
				}

			case maior-maior2 < 0.1:
				switch {
				case entidade == "SEGURANCA":
					switch d.Car {
					case m["SEGURANCA"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r2, `",`)
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r1, `",`)
					}
					break
				case entidade2 == "SEGURANCA":
					switch d.Car {
					case m["SEGURANCA"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r2, `",`)
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r1, `",`)
					}
					break
				case entidade == "CONSUMO":
					switch d.Car {
					case m["CONSUMO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r2, `",`)
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r1, `",`)
					}
					break
				case entidade2 == "CONSUMO":
					switch d.Car {
					case m["CONSUMO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r1, `",`)
						break
					}
					break
				case entidade == "DESEMPENHO":
					switch d.Car {
					case m["DESEMPENHO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r1, `",`)
						break
					}
					break
				case entidade2 == "DESEMPENHO":
					switch d.Car {
					case m["DESEMPENHO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r1, `",`)
						break
					}
					break
				case entidade == "MANUTENCAO":
					switch d.Car {
					case m["MANUTENCAO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r1, `",`)
						break
					}
					break
				case entidade2 == "MANUTENCAO":
					switch d.Car {
					case m["MANUTENCAO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r1, `",`)
						break
					}
					break
				case entidade == "CONFORTO":
					switch d.Car {
					case m["CONFORTO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r1, `",`)
						break
					}
					break
				case entidade2 == "CONFORTO":
					switch d.Car {
					case m["CONFORTO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r1, `",`)
						break
					}
					break
				case entidade == "DESIGN":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r1, `",`)
						break
					}
					break
				case entidade2 == "DESIGN":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r1, `",`)
						break
					}
					break
				case entidade == "ACESSORIOS":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r1, `",`)
						break
					}
					break
				case entidade2 == "ACESSORIOS":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r1, `",`)
						break
					}
					break
				}

			default:
				switch entidade {
				case "SEGURANCA":
					switch d.Car {
					case m["SEGURANCA"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["SEGURANCA"].r1, `",`)
						break
					}
					break
				case "CONSUMO":
					switch d.Car {
					case m["CONSUMO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONSUMO"].r1, `",`)
						break
					}
					break
				case "DESEMPENHO":
					switch d.Car {
					case m["DESEMPENHO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESEMPENHO"].r1, `",`)
						break
					}
					break
				case "MANUTENCAO":
					switch d.Car {
					case m["MANUTENCAO"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["MANUTENCAO"].r1, `",`)
						break
					}
					break
				case "CONFORTO":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["CONFORTO"].r1, `",`)
						break
					}
					break
				case "DESIGN":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["DESIGN"].r1, `",`)
						break
					}
					break
				case "ACESSORIOS":
					switch d.Car {
					case m["DESIGN"].r1:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r2, `",`)
						break
					default:
						recomendacao = gconcat.Build(recomendacao, `"`, m["ACESSORIOS"].r1, `",`)
						break
					}
					break
				}
			}

			switch math.Signbit(soma) {
			case false:
				return c.Status(200).JSON(`{"recommendation": "","entities": []}`)
			default:
				return c.Status(200).SendString(gconcat.Build(recomendacao, retorno))
			}

		default:
			return c.Status(500).SendString(responseErr1.Error())
		}

		return c.Status(200).SendString(`{"recommendation: "","entities": []}`)
	})
	app.Listen(":8081")
}
