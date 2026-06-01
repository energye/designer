// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package options

// web/index.html
const webIndexHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>ENERGY</title>
    <style>
        body {
            height: 100vh;
            margin: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            background: #f5f5f5;
        }
        button {
            padding: 10px 24px;
            border: 1px solid #ddd;
            border-radius: 8px;
            font-size: 16px;
            cursor: pointer;
        }
        button:hover {
            background: #eee;
        }
    </style>
</head>
<body>
<button id="count" onclick="count()">Count: -</button>
<script>
    let counter = 0;
    function count() {
        ipc.emit('counter:change', [++counter], function (r) {
            document.getElementById('count').innerHTML = "Count: " + r
        })
    }
</script>
</body>
</html>
`
