# FAQ

<!-- Unfortunately, GitHub does not support the markdown definition list -->

<dl>
<dt>
Why not CPython?
</dt>
<dd>
<p>
Integrating CPython requires the use of <a href="https://go.dev/wiki/cgo">cgo</a>. This makes it difficult for the average user to integrate the xk6-python extension into k6. The integration of CPython also means a runtime dependency. It is not enough to install the k6 executable, the CPython runtime must also be installed.
</p>
<p>
The <a href="https://github.com/google/starlark-go">starlark-go()</a> package is a pure go implementation of the <a href="https://github.com/google/starlark-go/blob/master/doc/spec.md">Starlark python dialect</a>. Its use does not require external dependencies, such as the installation of CPython. It doesn't even require the use of <a href="https://go.dev/wiki/cgo">cgo</a>. This enables the portability of the k6 executable binary and simplifies its use in the cloud.
</p>
</dd>

<dt>
Competing with JavaScript support?
</dt>
<dd>
<p>
No. Python language support complements k6 and its JavaScript language support. The target audience is different. The primary target audience is users biased towards the Python language.
</p>
<dt>
Will Python become a first-class citizen?
</dt>
<dd>
<p>
No. At least not at first. The level of Python support depends on how large a group of users wants to use k6 in Python.
</p>
<p>
Users who are committed to the Python language prefer to write tests in Python. Even if fewer k6 features are available in Python than in JavaScript.
</p>
</dd>
<dt>Competing with Python load testing tools?</dt>
<dd>
<p>
No. It is not intended to compete with the load testing tools that provide the entire Python ecosystem. The goal is to make k6 available to those who want to write load tests in Python, but do not need the entire Python ecosystem. To be honest, a load test usually doesn't even need the entire ecosystem.
</p>
<p>
The goal is also that those users who use Grafana products can choose a Grafana tool for load testing, even if they prefer the Python programming language.
</p>
</dd>
</dd>
</dl>
