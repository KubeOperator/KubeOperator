import { Component, Input, Output, OnInit, EventEmitter } from "@angular/core";
import { Subject } from "rxjs";
import { debounceTime } from 'rxjs/operators';

@Component({
  selector: "ko-filter",
  templateUrl: "./filter.component.html",
  styleUrls: ["./filter.component.scss"]
})
export class FilterComponent implements OnInit {
  placeHolder: string = "";
  filterTerms = new Subject<string>();
  isExpanded: boolean = false;

  @Output() private filterEvt = new EventEmitter<string>();
  @Output() private openFlag = new EventEmitter<boolean>();
  @Input() readonly: string = null;
  @Input() currentValue: string;
  @Input("filterPlaceholder")
  public set flPlaceholder(placeHolder: string) {
    this.placeHolder = placeHolder;
  }
  @Input() expandMode: boolean = false;
  @Input() withDivider: boolean = false;

  ngOnInit(): void {
    this.filterTerms
      .pipe(debounceTime(500))
      .subscribe(terms => {
        this.filterEvt.emit(terms);
      });
  }

  valueChange(): void {
    // Send out filter terms
    this.filterTerms.next(this.currentValue && this.currentValue.trim());
  }

  inputFocus(): void {
    this.openFlag.emit(this.isExpanded);
  }

  onClick(): void {
    // Only enabled when expandMode is set to false
    if (this.expandMode) {
      return;
    }
    this.isExpanded = !this.isExpanded;
    this.openFlag.emit(this.isExpanded);
  }

  public get isShowSearchBox(): boolean {
    return this.expandMode || (!this.expandMode && this.isExpanded);
  }
}
