const DESIGN_ICONS: Array<[string, string]> = [
  ['xd.adobe', 'https://upload.wikimedia.org/wikipedia/commons/c/c2/Adobe_XD_CC_icon.svg'],
  ['figma', 'https://upload.wikimedia.org/wikipedia/commons/3/33/Figma-logo.svg'],
  ['sketch', 'https://upload.wikimedia.org/wikipedia/commons/5/59/Sketch_Logo.svg'],
  ['zeplin', 'https://zeplin.io/img/icZeplin.svg'],
]

const LIBRARY_ICONS: Array<[string, string]> = [
  ['angular', 'https://coryrylan.com/assets/images/posts/types/angular.svg'],
  ['react', 'https://upload.wikimedia.org/wikipedia/commons/a/a7/React-icon.svg'],
  ['vue', 'https://upload.wikimedia.org/wikipedia/commons/9/95/Vue.js_Logo_2.svg'],
  ['svelte', 'https://upload.wikimedia.org/wikipedia/commons/1/1b/Svelte_Logo.svg'],
]

export function getIconDesignByUrl(url: string): string {
  return DESIGN_ICONS.find(([needle]) => url.includes(needle))?.[1] ?? ''
}

export function getIconLibraryByFilter(filter: string): string {
  return LIBRARY_ICONS.find(([needle]) => filter.includes(needle))?.[1] ?? ''
}

export function getVersionByFilter(filter: string): string {
  return filter.split('_')[1]
}
