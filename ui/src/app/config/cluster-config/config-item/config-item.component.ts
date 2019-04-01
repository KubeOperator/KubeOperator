import {Component, Input, OnInit} from '@angular/core';
import {ConfigBase} from '../config-base';
import {FormControl, FormGroup} from '@angular/forms';

@Component({
  selector: 'app-config-item',
  templateUrl: './config-item.component.html',
  styleUrls: ['./config-item.component.css']
})
export class ConfigItemComponent implements OnInit {

  constructor() {
  }

  ngOnInit() {
  }

  @Input() config: ConfigBase<any>;
  @Input() form: FormGroup;

  get isValid() {
    return this.form.controls[this.config.name].valid;
  }

}
