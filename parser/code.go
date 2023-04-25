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

// ToInput navigator.clipboard.writeText(text)
func (_this *Code) ToInput() string {
	return `<span class="input-group-text" id="basic-addon1">$name</span>
<input type="text" value="` + _this.ToString() + `" class="form-control" placeholder="code" aria-label="Code" aria-describedby="basic-addon1">
<button class="btn btn-outline-secondary" type="button" onclick="copyCode(this)">Copy</button>
`
}
func (_this *Code) ToCodeHtml() string {
	return `<div class="col"><pre><code class="language-go">` + _this.ToString() + `</code></pre></div>`
}

func (_this *Code) Value() []byte {
	return *_this
}
