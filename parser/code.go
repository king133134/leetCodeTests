package parser

type Code []byte

func (_this *Code) ToString() string {
	return string(*_this)
}

func (_this *Code) Append(b ...byte) {
	*_this = append(*_this, b...)
}

func (_this *Code) AppendBytes(b []byte) {
	*_this = append(*_this, b...)
}

func (_this *Code) ToCodeHtml() string {
	return `<div class="col"><pre><code class="language-go">` + _this.ToString() + `</code></pre></div>`
}

func (_this *Code) Value() []byte {
	return *_this
}
