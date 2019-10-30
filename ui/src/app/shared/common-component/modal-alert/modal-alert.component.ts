import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-modal-alert',
  templateUrl: './modal-alert.component.html',
  styleUrls: ['./modal-alert.component.css']
})
export class ModalAlertComponent implements OnInit {

  type: string;

  @Input()
  set message(message: string) {
    this._message = message;
  }
  get message(): string { return this._message; }
  _message: string;

  @Input()
  set tipShow(tipShow: boolean) {
    this._tipShow = tipShow;
  }
  get tipShow(): boolean { return this._tipShow; }
  _tipShow: boolean;

  @Input()
  set invalid(invalid: boolean) {
    this._invalid = invalid;
  }
  get invalid(): boolean { return this._invalid; }
  _invalid: boolean;

  constructor() { }

  ngOnInit() {
    this.tipShow = false;
  }

  closeTip() {
    this.tipShow = false;
  }
}
