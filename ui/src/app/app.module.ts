import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {BaseModule} from './base/base.module';
import {AccountModule} from './account/account.module';
import {InterceptorService} from './shared/interceptor.service';
import {HTTP_INTERCEPTORS} from '@angular/common/http';
import {PackageModule} from './package/package.module';
import {UserModule} from './user/user.module';
import {ClusterModule} from './cluster/cluster.module';
import {OverviewModule} from './overview/overview.module';
import {NodeModule} from './node/node.module';
import {LogModule} from './log/log.module';
import {HostModule} from './host/host.module';
import {DeployModule} from './deploy/deploy.module';
import {SettingModule} from './setting/setting.module';
import {CredentialModule} from './credential/credential.module';
import {RegionModule} from './region/region.module';
import {ZoneModule} from './zone/zone.module';
import {PlanModule} from './plan/plan.module';
import {F5BigIpModule} from './f5-big-ip/f5-big-ip.module';
import {ClusterHealthComponent} from './cluster-health/cluster-health.component';
import {DeployPlanComponent} from './deploy-plan/deploy-plan.component';
import {ApplicationComponent} from './application/application.component';
import {ClusterBackupComponent} from './cluster-backup/cluster-backup.component';
import {ClusterBackupModule} from './cluster-backup/cluster-backup.module';
import {SharedModule} from './shared/shared.module';
import {NfsModule} from './nfs/nfs.module';
import {StorageModule} from './storage/storage.module';
import {DashboardComponent} from './dashboard/dashboard.component';
import {SystemLogModule} from './system-log/system-log.module';
import {DnsModule} from './dns/dns.module';
import { ClusterStorageComponent } from './cluster-storage/cluster-storage.component';
import { ClusterEventComponent } from './cluster-event/cluster-event.component';
import { ClusterEventListComponent } from './cluster-event/cluster-event-list/cluster-event-list.component';
import { ClusterEventDetailComponent } from './cluster-event/cluster-event-detail/cluster-event-detail.component';
import { CephComponent } from './ceph/ceph.component';
import { CephListComponent } from './ceph/ceph-list/ceph-list.component';
import { CephCreateComponent } from './ceph/ceph-create/ceph-create.component';
import { NgCircleProgressModule } from 'ng-circle-progress';
import { ItemComponent } from './item/item.component';
import { ItemCreateComponent } from './item/item-create/item-create.component';
import { ItemListComponent } from './item/item-list/item-list.component';
import { ItemDetailComponent } from './item/item-detail/item-detail.component';
import { ItemMemberComponent } from './item-member/item-member.component';
import { ItemResourceComponent } from './item-resource/item-resource.component';
import { ItemResourceCreateComponent } from './item-resource/item-resource-create/item-resource-create.component';
import { ItemResourceListComponent } from './item-resource/item-resource-list/item-resource-list.component';
import { ItemResourcePipe } from './item-resource/item-resource.pipe';
import { ItemMemberListComponent } from './item-member/item-member-list/item-member-list.component';
import { ItemMemberCreateComponent } from './item-member/item-member-create/item-member-create.component';
import { MemberFilterPipe } from './item-member/member-filter.pipe';
import { PermissionNamePipe } from './item/permission-name.pipe';
import { ItemRolePipe } from './item-member/item-role.pipe';
import { MessageCenterComponent } from './message-center/message-center.component';
import { LocalMailComponent } from './message-center/local-mail/local-mail.component';
import { SubscribeComponent } from './message-center/subscribe/subscribe.component';
import { SubscribeConfigComponent } from './message-center/subscribe/subscribe-config/subscribe-config.component';
import { ReceiverComponent } from './message-center/receiver/receiver.component';
import { ClusterGradeComponent } from './cluster-grade/cluster-grade.component';
import { LocalMailDetailComponent } from './message-center/local-mail/local-mail-detail/local-mail-detail.component';

@NgModule({
  declarations: [
    AppComponent,
    DeployPlanComponent,
    ApplicationComponent,
    ClusterHealthComponent,
    ClusterBackupComponent,
    DashboardComponent,
    ClusterStorageComponent,
    ClusterEventComponent,
    ClusterEventListComponent,
    ClusterEventDetailComponent,
    CephComponent,
    CephListComponent,
    CephCreateComponent,
    ItemComponent,
    ItemCreateComponent,
    ItemListComponent,
    ItemDetailComponent,
    ItemMemberComponent,
    ItemResourceComponent,
    ItemResourceCreateComponent,
    ItemResourceListComponent,
    ItemResourcePipe,
    ItemMemberListComponent,
    ItemMemberCreateComponent,
    MemberFilterPipe,
    PermissionNamePipe,
    ItemRolePipe,
    MessageCenterComponent,
    LocalMailComponent,
    SubscribeComponent,
    SubscribeConfigComponent,
    ReceiverComponent,
    LocalMailDetailComponent,
    ClusterGradeComponent,
  ],
  imports: [
    CredentialModule,
    BrowserModule,
    BaseModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    AccountModule,
    PackageModule,
    UserModule,
    ClusterModule,
    DeployModule,
    OverviewModule,
    RegionModule,
    NodeModule,
    LogModule,
    HostModule,
    SettingModule,
    ZoneModule,
    PlanModule,
    F5BigIpModule,
    ClusterBackupModule,
    SharedModule,
    NfsModule,
    StorageModule,
    SystemLogModule,
    DnsModule,
       NgCircleProgressModule.forRoot({
      radius: 100,
      outerStrokeWidth: 16,
      innerStrokeWidth: 8,
      animationDuration: 300,
    })
  ],
  providers: [{provide: HTTP_INTERCEPTORS, useClass: InterceptorService, multi: true}],
  bootstrap: [AppComponent]
})
export class AppModule {
}
