package utils

import (
    "regexp"
    "strings"
)

func ReplaceStringByRegex(str, rule, replace string) (string, error) {
    reg, err := regexp.Compile(rule)
    if reg == nil || err != nil {
        return "", err
    }
    return reg.ReplaceAllString(str, replace), nil
}

func ConvertSymbolToUnderLine(s string) string {
    result, _ := ReplaceStringByRegex(s, "\\.", "_")
    result, _ = ReplaceStringByRegex(result, "-", "__")
    result, _ = ReplaceStringByRegex(result, "\\/", "___")
    result, _ = ReplaceStringByRegex(result, "=", "____")
    result, _ = ReplaceStringByRegex(result, ",", "_____")
    result, _ = ReplaceStringByRegex(result, ":", "a_a")
    result, _ = ReplaceStringByRegex(result, " ", "b_b")
    result, _ = ReplaceStringByRegex(result, "\\*", "c_c")
    if !strings.HasPrefix(result, "_") {
        return "_" + result
    } else {
        return result
    }
}

func GetFindAllData(pattern, matchDate string) []string {
    regex := regexp.MustCompile(pattern)
    return regex.FindAllString(matchDate, -1)
}
