import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';

@Component({
  selector: 'app-confirm-alert',
  templateUrl: './confirm-alert.component.html',
  styleUrls: ['./confirm-alert.component.css']
})
export class ConfirmAlertComponent implements OnInit {
  @Input() opened = false;
  @Output() openedChange = new EventEmitter();
  @Output() confirm = new EventEmitter();
  @Input() comment = '';
  @Input() title = '';


  constructor() {
  }

  ngOnInit() {
  }

  onCancel() {
    this.close();
  }

  setTitle(title: string) {
    this.title = title;
  }

  setComment(comment: string) {
    this.comment = comment;
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  onConfirm() {
    this.confirm.emit();
  }

}
