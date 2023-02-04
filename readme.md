# Shiki-github-graph

## Link

```shell
link = %host%/user/%name%
```

## Styles

```css
@media screen and (min-width: 1px) {
  .graph {
    display: none !important;
  }
  .activity .title {
    height: 150px;
  }
  .activity {
    overflow: clip;
  }
  .activity .title::after {
    content: url(%link%);
    right: 0px;
    position: absolute;
    display: block;
  }
}
```

## Example

default host: `http://shiki.mircloud.host/`.

example name: `POCCOMAXA`.

final link: `http://shiki.mircloud.host/user/POCCOMAXA`.

final style:

```css
@media screen and (min-width: 1px) {
  .graph {
    display: none !important;
  }
  .activity .title {
    height: 150px;
  }
  .activity {
    overflow: clip;
  }
  .activity .title::after {
    content: url(http://shiki.mircloud.host/user/POCCOMAXA);
    right: 0px;
    position: absolute;
    display: block;
  }
}
```
