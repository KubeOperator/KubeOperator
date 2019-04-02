import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NodeComponent} from './node.component';
import {CoreModule} from '../core/core.module';
import {NodeListComponent} from './node-list/node-list.component';
import {NodeService} from './node.service';
import { NodeCreateComponent } from './node-create/node-create.component';

@NgModule({
  declarations: [NodeComponent, NodeListComponent, NodeCreateComponent],
  imports: [
    CommonModule,
    CoreModule
  ], providers: [NodeService]
})
export class NodeModule {
}
