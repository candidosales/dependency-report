import { Pipe, PipeTransform } from '@angular/core';
import memo from 'memo-decorator';

@Pipe({
    name: 'iconDesign',
    standalone: true
})
export class IconDesignPipe implements PipeTransform {

  @memo()
  transform(url: string): string {
    return this.getIconDesignByUrl(url);
  }

  getIconDesignByUrl(url: string): string {
    if (url.includes('xd.adobe')) {
      return 'https://upload.wikimedia.org/wikipedia/commons/c/c2/Adobe_XD_CC_icon.svg';
    }
    if (url.includes('figma')) {
      return 'https://upload.wikimedia.org/wikipedia/commons/3/33/Figma-logo.svg';
    }
    if (url.includes('sketch')) {
      return 'https://upload.wikimedia.org/wikipedia/commons/5/59/Sketch_Logo.svg';
    }
    if (url.includes('zeplin')) {
      return 'https://zeplin.io/img/icZeplin.svg';
    }
    return '';
  }
}
