// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package update

import (
	"fmt"
	"net/http"
)

type IndexHandler struct {
}

func (i *IndexHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
	resp.Header().Set("Content-Type", "text/html")
	fmt.Fprint(resp, `<!DOCTYPE html>
<html>
	<head>
		<title>Updates</title>
		<script src="https://cdn.jsdelivr.net/npm/vue"></script>
	</head>
	<body>
		<div id="app">
			<h1>Updates</h1>
			<ul v-if="list.length > 0">
				<li v-for="element in list">
					{{ element.app }} {{ element.url }} {{ element.installedVersion }} < {{ element.latestVersion }}  
				</li>
			</ul>
			<p v-if="list.length == 0">Nothing to update</p>
		</div>
		<script>
			var demo = new Vue({
  				el: '#app',
				data: {
					list: [],
				},
				created: function () {
					this.fetchData()
				},
				methods: {
					fetchData: function () {
						var xhr = new XMLHttpRequest()
						var self = this
						xhr.open('GET', '/updates')
						xhr.onload = function () {
							self.list = JSON.parse(xhr.responseText)
						}
						xhr.send()
					},
				}
			})
		</script>
	</body>
</html>
`)
}
