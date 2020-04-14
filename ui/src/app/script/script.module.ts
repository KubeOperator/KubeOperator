import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ScriptComponent} from './script.component';
import {ScriptListComponent} from './script-list/script-list.component';
import {ScriptCreateComponent} from './script-create/script-create.component';
import {CoreModule} from '../core/core.module';
import {SharedModule} from '../shared/shared.module';
import {CodemirrorModule} from 'ng2-codemirror';


@NgModule({
  declarations: [ScriptComponent, ScriptListComponent, ScriptCreateComponent],
  imports: [
    CommonModule,
    CoreModule,
    SharedModule,
    CodemirrorModule
  ]
})
export class ScriptModule {
}
