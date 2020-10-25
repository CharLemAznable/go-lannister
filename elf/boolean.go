package elf

func ToBool(str string) bool {
    if "true" == str {
        return true
    }
    if "" == str {
        return false
    }
    switch len(str) {
    case 1:
        ch0 := str[0]
        if 'y' == ch0 || 'Y' == ch0 ||
            't' == ch0 || 'T' == ch0 ||
            '1' == ch0 {
            return true
        }
        if 'n' == ch0 || 'N' == ch0 ||
            'f' == ch0 || 'F' == ch0 ||
            '0' == ch0 {
            return false
        }
        break
    case 2:
        ch0 := str[0]
        ch1 := str[1]
        if ('o' == ch0 || 'O' == ch0) &&
            ('n' == ch1 || 'N' == ch1) {
            return true
        }
        if ('n' == ch0 || 'N' == ch0) &&
            ('o' == ch1 || 'O' == ch1) {
            return false
        }
        break
    case 3:
        ch0 := str[0]
        ch1 := str[1]
        ch2 := str[2]
        if ('y' == ch0 || 'Y' == ch0) &&
            ('e' == ch1 || 'E' == ch1) &&
            ('s' == ch2 || 'S' == ch2) {
            return true
        }
        if ('o' == ch0 || 'O' == ch0) &&
            ('f' == ch1 || 'F' == ch1) &&
            ('f' == ch2 || 'F' == ch2) {
            return false
        }
        break
    case 4:
        ch0 := str[0]
        ch1 := str[1]
        ch2 := str[2]
        ch3 := str[3]
        if ('t' == ch0 || 'T' == ch0) &&
            ('r' == ch1 || 'R' == ch1) &&
            ('u' == ch2 || 'U' == ch2) &&
            ('e' == ch3 || 'E' == ch3) {
            return true
        }
        break
    case 5:
        ch0 := str[0]
        ch1 := str[1]
        ch2 := str[2]
        ch3 := str[3]
        ch4 := str[4]
        if ('f' == ch0 || 'F' == ch0) &&
            ('a' == ch1 || 'A' == ch1) &&
            ('l' == ch2 || 'L' == ch2) &&
            ('s' == ch3 || 'S' == ch3) &&
            ('e' == ch4 || 'E' == ch4) {
            return false
        }
        break
    default:
        break
    }
    return false
}
