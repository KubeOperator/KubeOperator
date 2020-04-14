import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Script} from '../script';
import {ScriptService} from '../script.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';
import {NgForm} from '@angular/forms';
import * as CodeMirror from 'codemirror';


@Component({
  selector: 'app-script-create',
  templateUrl: './script-create.component.html',
  styleUrls: ['./script-create.component.css']
})
export class ScriptCreateComponent implements OnInit {

  opened = false;
  code = '';
  config = {};
  keys: string[] = [];
  item: Script = new Script();
  @Output() created = new EventEmitter<boolean>();
  @ViewChild('coder', {static: true}) coder;
  instance = null;
  @ViewChild('scriptForm', {static: true}) scriptFrom: NgForm;

  constructor(private service: ScriptService, private alert: CommonAlertService) {
  }

  ngOnInit() {
    this.config = {
      lineNumbers: true,
      theme: 'idea',
      mode: {name: 'shell'},
    };
    this.initEditor();
  }


  open() {
    this.scriptFrom.resetForm(this.item);
    this.opened = true;
  }

  initEditor() {
    this.instance = CodeMirror.fromTextArea(this.coder.nativeElement, this.config);
  }

  onSubmit() {
    this.service.create(this.item).subscribe(data => {
      this.opened = false;
      this.created.emit(true);
      this.alert.showAlert(`创建脚本成功`, AlertLevels.SUCCESS);
    }, error => {
      this.opened = false;
      this.created.emit(true);
      this.alert.showAlert(error, AlertLevels.ERROR);
    });
  }

  onTypeChange() {
    if (this.item) {
      if (this.item.type === 'other') {
        this.instance.setOption('mode', {name: 'text'});
      } else {
        this.instance.setOption('mode', {name: this.item.type});
      }
    }
  }

  onCancel() {
    this.opened = false;
  }

}
