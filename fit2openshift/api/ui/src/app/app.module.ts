import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';

import { AppComponent } from './app.component';
import { AppRoutingModule} from './app-routing.module';
import { CoreModule } from './core/core.module';
import { BaseModule } from './base/base.module';
import { SharedModule } from './shared/shared.module';
import { ProjectModule } from './project/project.module';
import { RoleModule } from './role/role.module';
import { AdhocModule } from './adhoc/adhoc.module';
import { PlaybookModule } from './playbook/playbook.module';
import { HostModule } from './host/host.module';
import { GroupModule } from './group/group.module';
import { CeleryModule } from './celery/celery.module';


@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    RouterModule,
    CoreModule,
    BaseModule,
    SharedModule,
    AppRoutingModule,
    ProjectModule,
    PlaybookModule,
    RoleModule,
    AdhocModule,
    CeleryModule,
    HostModule,
    GroupModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
