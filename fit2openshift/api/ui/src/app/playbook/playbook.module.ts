import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CoreModule } from '../core/core.module';
import { SharedModule } from '../shared/shared.module';

import { CeleryModule } from '../celery/celery.module';
import { PlaybookComponent } from './playbook.component';
import { PlaybookListComponent } from './playbook-list/playbook-list.component';
import { PlaybookExecuteComponent } from './playbook-execute/playbook-execute.component';
import { PlaybookCreateComponent } from './playbook-create/playbook-create.component';

@NgModule({
  imports: [
    CommonModule,
    SharedModule,
    CoreModule,
    CeleryModule,
  ],
  declarations: [PlaybookListComponent, PlaybookExecuteComponent, PlaybookComponent, PlaybookCreateComponent]
})
export class PlaybookModule { }
