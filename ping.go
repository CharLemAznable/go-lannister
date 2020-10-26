package lannister

import (
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/mvc"
)

type PingController struct{}

func (c *PingController) BeforeActivation(b mvc.BeforeActivation) {
    b.Handle(GetMapping("/", "Index"))
    b.Handle(GetMapping("/welcome", "Welcome"))
}

func (c *PingController) Index(ctx iris.Context) {
    ctx.Redirect(ctx.Path() + "/welcome")
}

func (c *PingController) Welcome() string {
    return "Lannisters Always pay their debts."
}

func init() {
    RegisterController("lannister.PingController", &PingController{})
}
