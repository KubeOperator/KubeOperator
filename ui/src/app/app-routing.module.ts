import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {LoginComponent} from './login/login.component';
import {LayoutComponent} from './layout/layout.component';
import {ClusterComponent} from './business/cluster/cluster.component';
import {ClusterDetailComponent} from './business/cluster/cluster-detail/cluster-detail.component';
import {ClusterRoutingResolverService} from './business/cluster/cluster-routing-resolver.service';
import {OverviewComponent} from './business/cluster/cluster-detail/overview/overview.component';
import {SettingComponent} from './business/setting/setting.component';
import {CredentialComponent} from './business/setting/credential/credential.component';
import {HostComponent} from './business/host/host.component';
import {NodeComponent} from './business/cluster/cluster-detail/node/node.component';
import {NamespaceComponent} from './business/cluster/cluster-detail/namespace/namespace.component';
import {StorageComponent} from './business/cluster/cluster-detail/storage/storage.component';
import {PersistentVolumeComponent} from './business/cluster/cluster-detail/storage/persistent-volume/persistent-volume.component';
import {PersistentVolumeClaimComponent} from './business/cluster/cluster-detail/storage/persistent-volume-claim/persistent-volume-claim.component';
import {UserComponent} from './business/user/user.component';
import {AuthUserService} from './shared/auth/auth-user.service';
import {LoggingComponent} from './business/cluster/cluster-detail/logging/logging.component';
import {MonitorComponent} from './business/cluster/cluster-detail/monitor/monitor.component';
import {StorageClassComponent} from './business/cluster/cluster-detail/storage/storage-class/storage-class.component';
import {RegionComponent} from './business/deploy-plan/region/region.component';
import {DeployPlanComponent} from './business/deploy-plan/deploy-plan.component';
import {ZoneComponent} from './business/deploy-plan/zone/zone.component';
import {PlanComponent} from './business/deploy-plan/plan/plan.component';
import {StorageProvisionerComponent} from './business/cluster/cluster-detail/storage/storage-provisioner/storage-provisioner.component';
import {RepositoryComponent} from './business/cluster/cluster-detail/repository/repository.component';
import {ChartmuseumComponent} from './business/cluster/cluster-detail/repository/chartmuseum/chartmuseum.component';
import {RegistryComponent} from './business/cluster/cluster-detail/repository/registry/registry.component';
import {ToolsComponent} from './business/cluster/cluster-detail/tools/tools.component';
import {SystemComponent} from './business/setting/system/system.component';
import {ProjectComponent} from './business/project/project.component';
import {ProjectDetailComponent} from './business/project/project-detail/project-detail.component';
import {ProjectRoutingResolverService} from './business/project/project-routing-resolver.service';
import {ProjectResourceComponent} from './business/project/project-resource/project-resource.component';
import {ProjectMemberComponent} from './business/project/project-member/project-member.component';
import {LogComponent} from './business/cluster/cluster-detail/log/log.component';
import {BackupAccountComponent} from './business/setting/backup-account/backup-account.component';
import {BackupComponent} from './business/cluster/cluster-detail/backup/backup.component';
import {LicenseComponent} from './business/setting/license/license.component';
import {SecurityComponent} from './business/cluster/cluster-detail/security/security.component';
import {LdapComponent} from './business/setting/ldap/ldap.component';
import {ManifestComponent} from './business/manifest/manifest.component';
import {ThemeComponent} from './business/setting/theme/theme.component';
import {EventComponent} from './business/cluster/cluster-detail/event/event.component';
import {MessageCenterComponent} from './business/message-center/message-center.component';
import {UserReceiverComponent} from './business/message-center/user-receiver/user-receiver.component';
import {UserSubscribeComponent} from './business/message-center/user-subscribe/user-subscribe.component';
import {MailboxComponent} from './business/message-center/mailbox/mailbox.component';
import {ClusterLoggerComponent} from './business/cluster/cluster-logger/cluster-logger.component';
import {MessageComponent} from './business/setting/message/message.component';
import {VmConfigComponent} from './business/deploy-plan/vm-config/vm-config.component';
import {ClusterGradeComponent} from './business/cluster/cluster-detail/cluster-grade/cluster-grade.component';
import {F5Component} from './business/cluster/cluster-detail/f5/f5.component';
import {BusinessResolverService} from './shared/service/business-resolver.service';
import {AdminAuthService} from './shared/auth/admin-auth.service';
import {EmailComponent} from './business/setting/email/email.component';
import {SystemLogComponent} from './business/system-log/system-log.component';
import {MultiClusterComponent} from "./business/multi-cluster/multi-cluster.component";

const routes: Routes = [
    {path: 'login', component: LoginComponent},
    {path: 'logger', component: ClusterLoggerComponent},
    {
        path: '',
        component: LayoutComponent,
        canActivate: [AuthUserService],
        canActivateChild: [AuthUserService],
        resolve: {hasLicense: BusinessResolverService},
        children: [
            {path: '', redirectTo: 'projects', pathMatch: 'full'},
            {
                path: 'projects',
                component: ProjectComponent,
            },
            {
                path: 'projects/:name',
                component: ProjectDetailComponent,
                resolve: {project: ProjectRoutingResolverService},
                children: [
                    {path: '', redirectTo: 'clusters', pathMatch: 'full'},
                    {path: 'clusters', component: ClusterComponent},
                    {path: 'resources', component: ProjectResourceComponent},
                    {path: 'members', component: ProjectMemberComponent},
                ]
            },
            {
                path: 'projects/:projectName/clusters/:name',
                component: ClusterDetailComponent,
                resolve: {cluster: ClusterRoutingResolverService},
                children: [
                    {path: '', redirectTo: 'overview', pathMatch: 'full'},
                    {path: 'overview', component: OverviewComponent},
                    {path: 'nodes', component: NodeComponent},
                    {path: 'namespaces', component: NamespaceComponent},
                    {path: 'events', component: EventComponent},
                    {
                        path: 'storages',
                        component: StorageComponent,
                        children: [
                            {path: '', redirectTo: 'pv', pathMatch: 'full'},
                            {path: 'pv', component: PersistentVolumeComponent},
                            {path: 'pvc', component: PersistentVolumeClaimComponent},
                            {path: 'sc', component: StorageClassComponent},
                            {path: 'provisioner', component: StorageProvisionerComponent},
                        ],
                    },
                    {path: 'logging', component: LoggingComponent},
                    {path: 'monitor', component: MonitorComponent},
                    {path: 'security', component: SecurityComponent},
                    {
                        path: 'repository',
                        component: RepositoryComponent,
                        children: [
                            {path: '', redirectTo: 'chartmuseum', pathMatch: 'full'},
                            {path: 'chartmuseum', component: ChartmuseumComponent},
                            {path: 'registry', component: RegistryComponent}
                        ]
                    },
                    {path: 'tool', component: ToolsComponent},
                    {path: 'backup', component: BackupComponent},
                    {path: 'logs', component: LogComponent},
                    {path: 'grade', component: ClusterGradeComponent},
                    {path: 'f5', component: F5Component}
                ],
            },
            {
                path: 'hosts',
                component: HostComponent,
                canActivate: [AdminAuthService]
            },
            {
                path: 'muticluster',
                component: MultiClusterComponent,
                canActivate: [AdminAuthService]
            },
            {
                path: 'setting',
                component: SettingComponent,
                canActivate: [AdminAuthService],
                canActivateChild: [AdminAuthService],
                children: [
                    {path: '', redirectTo: 'system', pathMatch: 'full'},
                    {path: 'system', component: SystemComponent},
                    {path: 'credential', component: CredentialComponent},
                    {path: 'backupAccounts', component: BackupAccountComponent},
                    {path: 'email', component: EmailComponent},
                    {path: 'license', component: LicenseComponent},
                    {path: 'ldap', component: LdapComponent},
                    {path: 'theme', component: ThemeComponent},
                    {path: 'message', component: MessageComponent},
                ]
            },
            {
                path: 'deploy',
                component: DeployPlanComponent,
                canActivate: [AdminAuthService],
                children: [
                    {path: '', redirectTo: 'region', pathMatch: 'full'},
                    {path: 'region', component: RegionComponent},
                    {path: 'zone', component: ZoneComponent},
                    {path: 'plan', component: PlanComponent},
                    {path: 'vmConfig', component: VmConfigComponent}
                ]
            }, {
                path: 'manifests',
                component: ManifestComponent,
                canActivate: [AdminAuthService],
            },
            {
                path: 'users',
                component: UserComponent,
                canActivate: [AdminAuthService],
            },
            {
                path: 'system-log',
                component: SystemLogComponent,
                canActivate: [AdminAuthService],
            },
            {
                path: 'message',
                component: MessageCenterComponent,
                children: [
                    {path: '', redirectTo: 'mailbox', pathMatch: 'full'},
                    {path: 'userReceiver', component: UserReceiverComponent},
                    {path: 'subscribe', component: UserSubscribeComponent},
                    {path: 'mailbox', component: MailboxComponent}
                ]
            }
        ]
    }
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule {
}
