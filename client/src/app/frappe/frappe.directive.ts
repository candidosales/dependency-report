import { Directive, ElementRef, Input, Output, EventEmitter, OnChanges } from '@angular/core';
declare var frappe: any;

@Directive({
  selector: '[appFrappe]'
})
export class FrappeDirective implements OnChanges {

  @Input() title: string;
  @Input() data: any;
  @Input() type: string;
  @Input() height: number;

  @Output() frappe: EventEmitter<any> = new EventEmitter();

  constructor(
    private el: ElementRef
  ) {
    this.type = this.type || 'bar';
    this.height = this.height || 250;
  }

  ngOnChanges() {
    const chart = new frappe.Chart(
      this.el.nativeElement,
      {
        title: this.title,
        data: this.data,
        type: this.type, // or 'line', 'scatter', 'pie', 'percentage'
        height: this.height
      }
    );
    this.frappe.emit(chart);
  }
}
