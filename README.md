# LLMT - Large Language Model Transformer

Transforms all files in specified directory using a large language model. You can specify conditions for file
transformation along with a prompt, which LLM to use and which model should be used . The transformed files are saved in
a new directory. The file structure is preserved.


## Installation

### Command line utility

#### Docker

```shell
docker run -v $(pwd):/data github.com/blwsh/llmt analyze
```

#### Release

Download the latest release from the [releases page](https://github.com/blwsh/llmt/releases). Extract the archive and run the binary.

## Command line usage

```shell
llmt \ --config <config_file> \    # optional parameter, default is config.yaml in current directory
      analyze ./myProject ../docs  # analyzes files in ./myProject and outputs them as markdown in ../docs (maintains file structure)
```

Example `config.yaml` file:

```yaml
#$schema: https://raw.githubusercontent.com/blwsh/llmt/main/schema.json
version: "0.1"

analyzers:
- prompt: What is your name?
  analyzer: openai
  model: gpt-4o-mini
  regex: ^.+\.php$
  not_in:
    - vendor
```

See [schema.json](schema.json) for the full config schema.

## Go usage

You can find a comprehensive list of [examples](examples) here. Below is a simple example which has similar behaviour to the command line analyze command.

```go
package main

import (
	"context"

	"github.com/blwsh/llmt/pkg/file_analyzer/openai"
	"github.com/blwsh/llmt/pkg/project_analyzer"
)

func main() {
	ctx := context.Background()

	project_analyzer.New().
		AnalyzeProject(ctx, "myProject", "../docs", []project_analyzer.FileAnalyzer{
			{
				Prompt:        "document this files behaviour",
				Analyzer:      openai.New(openAIToken, "gpt-4o-mini"),
				Condition:     myFancyConditionFunc,
				ResultHandler: myDocsWriterFunc,
			},
		})
}
```

With `Condition` and `ResultHandler` you're able to filter out which files should be processed and how the result should be processed.

To see a more complete example of the above snippet, see [examples/overview](examples/overview/main.go) directory.


