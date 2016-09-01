if exists("b:current_syntax")
    finish
endif

syntax keyword legionFunction run copy ex call set

syntax match legionComment "\v#.*$"

syntax match legionOperator "{.*}"

highlight link legionComment Comment
highlight link legionFunction Function
highlight link legionOperator Operator

let b:current_syntax = "legion"
