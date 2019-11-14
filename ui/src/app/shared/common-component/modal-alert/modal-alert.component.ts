import {Component, Input, OnInit, Output, EventEmitter} from '@angular/core';

@Component({
  selector: 'app-modal-alert',
  templateUrl: './modal-alert.component.html',
  styleUrls: ['./modal-alert.component.css']
})
export class ModalAlertComponent implements OnInit {


  @Input() public tipShow: boolean;
  @Output() public tipShowChange = new EventEmitter();
  @Input() public message: string;
  @Output() public messageChange = new EventEmitter();
  @Input() public invalid: boolean;
  @Output() public invalidChange = new EventEmitter();

  constructor() {
  }

  ngOnInit() {
    this.tipShow = false;
  }

  closeTip() {
    this.tipShow = false;
  }

  showTip(invalid, message) {
    this.tipShow = true;
    this.message = message;
    this.invalid = invalid;
  }
}
