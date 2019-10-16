import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NfsComponent} from './nfs.component';
import {NfsListComponent} from './nfs-list/nfs-list.component';
import {NfsCreateComponent} from './nfs-create/nfs-create.component';
import {CoreModule} from '../core/core.module';


@NgModule({
  declarations: [NfsComponent, NfsListComponent, NfsCreateComponent],
  imports: [
    CommonModule,
    CoreModule
  ]
})
export class NfsModule {
}
