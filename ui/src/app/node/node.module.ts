import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NodeComponent} from './node.component';
import {CoreModule} from '../core/core.module';
import {NodeListComponent} from './node-list/node-list.component';
import {NodeService} from './node.service';
import { NodeDetailComponent } from './node-detail/node-detail.component';
import {NgCircleProgressModule} from 'ng-circle-progress';
import {SharedModule} from "../shared/shared.module";

@NgModule({
  declarations: [NodeComponent, NodeListComponent, NodeDetailComponent],
  imports: [
    CommonModule,
    CoreModule,
    NgCircleProgressModule,
    SharedModule
  ], providers: [NodeService]
})
export class NodeModule {
}
