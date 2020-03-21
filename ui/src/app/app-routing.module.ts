import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {ShellComponent} from './base/shell/shell.component';
import {NotFoundComponent} from './shared/not-found/not-found.component';
import {SignInComponent} from './account/sign-in/sign-in.component';
import {AuthUserActiveService} from './shared/route/auth-user-active.service';
import {PackageComponent} from './package/package.component';
import {UserComponent} from './user/user.component';
import {ClusterComponent} from './cluster/cluster.component';
import {ClusterDetailComponent} from './cluster/cluster-detail/cluster-detail.component';
import {OverviewComponent} from './overview/overview.component';
import {NodeComponent} from './node/node.component';
import {LogComponent} from './log/log.component';
import {ClusterRoutingResolverService} from './cluster/cluster-routing-resolver.service';
import {HostComponent} from './host/host.component';
import {DeployComponent} from './deploy/deploy.component';
import {SettingComponent} from './setting/setting.component';
import {SystemSettingComponent} from './setting/system-setting/system-setting.component';
import {CredentialComponent} from './credential/credential.component';
import {RegionComponent} from './region/region.component';
import {ZoneComponent} from './zone/zone.component';
import {PlanComponent} from './plan/plan.component';
import {F5BigIpComponent} from './f5-big-ip/f5-big-ip.component';
import {DeployPlanComponent} from './deploy-plan/deploy-plan.component';
import {ClusterHealthComponent} from './cluster-health/cluster-health.component';
import {ApplicationComponent} from './application/application.component';
import {ClusterBackupComponent} from './cluster-backup/cluster-backup.component';
import {BackupStorageSettingComponent} from './setting/backup-storage-setting/backup-storage-setting.component';
import {NfsComponent} from './nfs/nfs.component';
import {StorageComponent} from './storage/storage.component';
import {DashboardComponent} from './dashboard/dashboard.component';
import {SystemLogComponent} from './system-log/system-log.component';
import {DnsComponent} from './dns/dns.component';
import {ClusterStorageComponent} from './cluster-storage/cluster-storage.component';
import {ClusterEventComponent} from './cluster-event/cluster-event.component';
import {CephComponent} from './ceph/ceph.component';
import {ItemComponent} from './item/item.component';
import {ItemDetailComponent} from './item/item-detail/item-detail.component';
import {ItemRoutingResolverService} from './item/item-routing-resolver.service';
import {ItemMemberComponent} from './item-member/item-member.component';
import {ItemResourceComponent} from './item-resource/item-resource.component';
import {LdapComponent} from './setting/ldap/ldap.component';
import {NotificationComponent} from './setting/notification/notification.component';
import {MessageCenterComponent} from './message-center/message-center.component';
import {LocalMailComponent} from './message-center/local-mail/local-mail.component';
import {DescribeComponent} from './overview/describe/describe.component';
import {SubscribeComponent} from './message-center/subscribe/subscribe.component';
import {ReceiverComponent} from './message-center/receiver/receiver.component';
import {ClusterGradeComponent} from './cluster-grade/cluster-grade.component';

const routes: Routes = [
  {path: 'sign-in', component: SignInComponent},
  {
    path: '',
    component: ShellComponent,
    canActivate: [AuthUserActiveService],
    canActivateChild: [AuthUserActiveService],
    children: [
      {path: '', redirectTo: 'dashboard', pathMatch: 'full'},
      {path: 'dashboard', component: DashboardComponent},
      {path: 'cluster', component: ClusterComponent},
      {path: 'item', component: ItemComponent},
      {path: 'package', component: PackageComponent},
      {path: 'user', component: UserComponent},
      {path: 'host', component: HostComponent},
      {path: 'user', component: UserComponent},
      {
        path: 'storage',
        component: StorageComponent,
        children: [
          {path: '', redirectTo: 'nfs', pathMatch: 'full'},
          {path: 'nfs', component: NfsComponent},
          {path: 'ceph', component: CephComponent},
        ]
      },
      {
        path: 'plan',
        component: DeployPlanComponent,
        children: [
          {path: '', redirectTo: 'region', pathMatch: 'full'},
          {path: 'region', component: RegionComponent},
          {path: 'zone', component: ZoneComponent},
          {path: 'plan', component: PlanComponent}
        ]
      },
      {
        path: 'setting',
        component: SettingComponent,
        children: [
          {path: '', redirectTo: 'system', pathMatch: 'full'},
          {path: 'system', component: SystemSettingComponent},
          {path: 'credential', component: CredentialComponent},
          {path: 'backup-storage', component: BackupStorageSettingComponent},
          {path: 'ldap', component: LdapComponent},
          {path: 'notification', component: NotificationComponent}
        ]
      },
      {
        path: 'cluster/:name',
        component: ClusterDetailComponent,
        resolve: {cluster: ClusterRoutingResolverService},
        children: [
          {path: '', redirectTo: 'overview', pathMatch: 'full'},
          {path: 'overview', component: OverviewComponent},
          {path: 'node', component: NodeComponent},
          {path: 'deploy', component: DeployComponent},
          {path: 'log', component: LogComponent},
          {path: 'apps', component: ApplicationComponent},
          {path: 'health', component: ClusterHealthComponent},
          {path: 'event', component: ClusterEventComponent},
          {path: 'backup', component: ClusterBackupComponent},
          {path: 'grade', component: ClusterGradeComponent},
          {path: 'big-ip', component: F5BigIpComponent},
          {path: 'cluster-storage', component: ClusterStorageComponent}
        ]
      },
      {
        path: 'item/:itemName',
        component: ItemDetailComponent,
        resolve: {item: ItemRoutingResolverService},
        children: [
          {path: '', redirectTo: 'cluster', pathMatch: 'full'},
          {
            path: 'cluster', component: ClusterComponent
          },
          {path: 'members', component: ItemMemberComponent},
          {path: 'resource', component: ItemResourceComponent},
        ]
      },
      {
        path: 'system/log', component: SystemLogComponent
      },
      {
        path: 'messageCenter',
        component: MessageCenterComponent,
        children: [
          {path: '', redirectTo: 'localMail', pathMatch: 'full'},
          {path: 'localMail', component: LocalMailComponent},
          {path: 'subscribe', component: SubscribeComponent},
          {path: 'receiver', component: ReceiverComponent},
        ]
      }
    ]
  },
  {path: '**', component: NotFoundComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
