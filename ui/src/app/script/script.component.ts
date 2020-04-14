import {Component, OnInit, ViewChild} from '@angular/core';
import {ScriptListComponent} from './script-list/script-list.component';
import {ScriptCreateComponent} from './script-create/script-create.component';

@Component({
  selector: 'app-script',
  templateUrl: './script.component.html',
  styleUrls: ['./script.component.css']
})
export class ScriptComponent implements OnInit {

  constructor() {
  }


  @ViewChild(ScriptListComponent, {static: true})
  list: ScriptListComponent;

  @ViewChild(ScriptCreateComponent, {static: true})
  creation: ScriptCreateComponent;

  ngOnInit() {
  }

  openCreate() {
    this.creation.open();
  }

  afterCreate() {
    this.list.list();
  }

}
