package calculate

import (
    "github.com/dengsgo/math-engine/engine"
    log "github.com/sirupsen/logrus"
    "regexp"
    "strings"
)

func Exec() {
    var exp string
    exp = "Cpu_used/(Cpu_used+Cpu_free)*100"
    reg := regexp.MustCompile(`[A-Za-z]*._.[A-Za-z]*`)
    if reg == nil {
        log.Errorf("regexp MustCompile err: %v", reg)
        return
    }
    result := reg.FindAllStringSubmatch(exp, -1)
    log.Infof("Exec result: %v", result)
    for _, v := range result {
        if v[0] == "Cpu_used" {
            exp = strings.ReplaceAll(exp, v[0], "45")
        }
        if v[0] == "Cpu_free" {
            exp = strings.ReplaceAll(exp, v[0], "55")
        }
    }
    r, err := engine.ParseAndExec(exp)
    if err != nil {
        log.Errorf("engine ParseAndExec err: %v", err)
    }
    log.Infof("%s = %v", exp, r)

    toks, err := engine.Parse(exp)
    if err != nil {
        log.Errorf("engine Parse ERROR: %v", err.Error())
        return
    }

    ast := engine.NewAST(toks, exp)
    if ast.Err != nil {
        log.Errorf("engine NewAST ERROR: %v", ast.Err.Error())
        return
    }

    ar := ast.ParseExpression()
    if ast.Err != nil {
        log.Errorf("ast ParseExpression ERROR: %v", ast.Err.Error())
        return
    }

    r = engine.ExprASTResult(ar)
    log.Infof("%s = %v\n", exp, r)
}
