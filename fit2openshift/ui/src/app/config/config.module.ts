import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ConfigComponent} from './config.component';
import {NodeConfigComponent} from './node-config/node-config.component';
import {ClusterConfigComponent} from './cluster-config/cluster-config.component';
import {CoreModule} from '../core/core.module';
import {ConfigService} from './config.service';

@NgModule({
  declarations: [ConfigComponent, NodeConfigComponent, ClusterConfigComponent],
  imports: [
    CommonModule,
    CoreModule
  ],
  providers: [ConfigService]
})
export class ConfigModule {
}
