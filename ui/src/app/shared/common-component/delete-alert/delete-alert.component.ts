import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';

@Component({
  selector: 'app-delete-alert',
  templateUrl: './delete-alert.component.html',
  styleUrls: ['./delete-alert.component.css']
})
export class DeleteAlertComponent implements OnInit {
  @Input() opened = false;
  @Output() openedChange = new EventEmitter();
  @Output() confirm = new EventEmitter();
  @Input() resourceTypeName = '未知';
  @Input() msg = '无';
  @Input() items: any[] = [];

  constructor() {
  }

  ngOnInit() {
  }

  close() {
    this.opened = false;
    this.openedChange.emit(this.opened);
  }

  confirmDelete() {
    this.confirm.emit();
  }
}
