# render CloudFormation to Json

Extracts the property part of the specified resource from the specified CloudFormation file and creates a json.

## Inputs

### `name`

**Required** Specify the resource name to be extracted.

### `json-params`

Specify the CloudFormation parameter in json format.

```js:example
'{"EnvType":"development"}'
```

### `cfm-definitoion`

**Required** Specify the path of the CloudFormation file.

## Outputs

### `json-file`

Temporary file names that can be used as-is

## Example usage
```yaml
steps:
  - name: Render cfm to task definition
    id: render-json
    uses: goccha/render-cfm-to-json@v0.0.1
    with:
      name: ecsTask
      json-params: '{"EnvType":"development"}'
      cfm-definition: deployments/ecs-task.yaml

  - name: Render new task definition
    id: render-container
    uses: aws-actions/amazon-ecs-render-task-definition@v1
    with:
      task-definition: ${{ steps.render-json.outputs.json-file }}
      container-name: api
      image: gchr.io/{owner}/{image-name}:{tag}

```