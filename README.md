# LLMT - Large Language Model Transformer

Transforms all files in specified directory using a large language model. You can specify conditions for file
transformation along with a prompt, which LLM to use and which model should be used . The transformed files are saved in
a new directory. The file structure is preserved.

## Usage

### Command line usage

#### Docker

```shell
docker run -v $(pwd):/data github.com/blwsh/llmt analyze
```

> [!NOTE]  
> When using `openai` analyzer you'll need to provide an API key. You can do this by setting the `OPENAI_TOKEN`
> environment variable.

#### Release

Download the latest release from the [releases page](https://github.com/blwsh/llmt/releases). Extract the archive and
run the binary.

```shell
llmt \ --config <config_file> \    # optional parameter, default is config.yaml in current directory
      analyze ./myProject ../docs  # analyzes files in ./myProject and outputs them as markdown in ../docs (maintains file structure)
```

### Configuration

Example `config.yaml` file:

```yaml
#$schema: https://raw.githubusercontent.com/blwsh/llmt/main/schema.json
version: "0.1"

analyzers:
  - prompt: Write docs for this file
    analyzer: openai
    model: gpt-4o-mini
    regex: ^.+\.php$
    not_in:
      - vendor
```

#### Analyzer configuration

| Field    | Description                                                                                                 | Required | Default |
|----------|-------------------------------------------------------------------------------------------------------------|----------|---------|
| prompt   | The prompt to use for the language model                                                                    | yes      |         |
| analyzer | Specifies which llm to use. See [Available analyzers](#available-analyzers) for the full list of analyzers. | yes      |         |
| model    | The model to use for the language model. If you use a fine tuned openai model, you set its name here.       | yes      |         |
| regex    | A regex to match the file path.                                                                             | no       |         |
| not_in   | A list of directories to exclude from the analysis.                                                         | no       |         |
| in       | A list of directories to include in the analysis. Note: not_in takes precedence over `in`.                  | no       |         |

See [schema.json](schema.json) for the full config schema.

#### Available analyzers

| Analyzer | Description                                                                                                                    |
|----------|--------------------------------------------------------------------------------------------------------------------------------|
| openai   | Uses the OpenAI API to transform the files. You need to provide an API key via the `OPENAI_TOKEN` environment variable.        |
| ollama   | Uses the OLLAMA API to transform the files. You can override the ollama url by setting the `OLLAMA_HOST` environment variable. |

## Go usage

You can find a comprehensive list of [examples](examples) here. Below is a simple example which has similar behaviour to
the command line analyze command.

<details>
  <summary>Click to expand!</summary>

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
    
    > [!TIP]
    > To see a more complete example of the above snippet, see [examples/overview](examples/overview/main.go) directory.

</details>
