if exists("b:current_syntax")
    finish
endif

syntax keyword legionFunction cmd copy config echo debug halt include

syntax match legionComment "\v#.*$"

syntax match legionOperator "{[^}]*}"

highlight link legionComment Comment
highlight link legionFunction Function
highlight link legionOperator Operator

let b:current_syntax = "legion"
