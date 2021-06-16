if exists("b:current_syntax")
    finish
endif

syntax keyword legionFunction CMD COPY CONFIG ECHO DEBUG HALT INCLUDE

syntax match legionComment "\v#.*$"

syntax match legionOperator "{[^}]*}"

highlight link legionComment Comment
highlight link legionFunction Function
highlight link legionOperator Operator

let b:current_syntax = "legion"
