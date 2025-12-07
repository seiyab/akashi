package doc

import "strings"

func Inline(text string) inline {
	return inline{onRight: true, onLeft: true, text: text}
}

type inline struct {
	onLeft  bool
	onRight bool
	text    string
}

func (l inline) print(depth int) virtualLines {
	var marker string = "  "
	side := lineNeutral
	isLeft := l.onLeft && !l.onRight
	isRight := l.onRight && !l.onLeft
	if isLeft {
		marker = "- "
		side = lineLeft
	} else if isRight {
		marker = "+ "
		side = lineRight
	}
	return virtualLines{
		{
			isDiff: isLeft || isRight,
			side:   side,
			text:   marker + strings.Repeat(" ", depth*indentSize) + l.text,
		},
	}
}

func (l inline) Left() Doc {
	l.onLeft = true
	l.onRight = false
	return l
}

func (l inline) Right() Doc {
	l.onLeft = false
	l.onRight = true
	return l
}

func (l inline) AddPrefix(prefix string) Doc {
	l.text = prefix + l.text
	return l
}

func (l inline) AddSuffix(suffix string) Doc {
	l.text = l.text + suffix
	return l
}
