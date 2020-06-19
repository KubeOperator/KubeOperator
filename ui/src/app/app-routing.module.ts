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
import {WorkloadComponent} from './business/cluster/cluster-detail/workload/workload.component';
import {DeploymentComponent} from './business/cluster/cluster-detail/workload/deployment/deployment.component';
import {StatefulSetComponent} from './business/cluster/cluster-detail/workload/stateful-set/stateful-set.component';
import {DaemonSetComponent} from './business/cluster/cluster-detail/workload/daemon-set/daemon-set.component';
import {JobComponent} from './business/cluster/cluster-detail/workload/job/job.component';
import {CornJobComponent} from './business/cluster/cluster-detail/workload/corn-job/corn-job.component';
import {ServiceComponent} from './business/cluster/cluster-detail/service/service.component';
import {IngressComponent} from './business/cluster/cluster-detail/ingress/ingress.component';
import {UserComponent} from './business/user/user.component';
import {AuthUserService} from './shared/auth/auth-user.service';
import {ConfigMapComponent} from './business/cluster/cluster-detail/config/config-map/config-map.component';
import {SecretComponent} from './business/cluster/cluster-detail/config/secret/secret.component';
import {ConfigComponent} from './business/cluster/cluster-detail/config/config.component';
import {LoggingComponent} from './business/cluster/cluster-detail/logging/logging.component';
import {TaskComponent} from './business/cluster/cluster-detail/task/task.component';
import {MonitorComponent} from './business/cluster/cluster-detail/monitor/monitor.component';

const routes: Routes = [
    {path: 'login', component: LoginComponent},
    {
        path: '',
        component: LayoutComponent,
        canActivate: [AuthUserService],
        canActivateChild: [AuthUserService],
        children: [
            {path: '', redirectTo: 'clusters', pathMatch: 'full'},
            {
                path: 'clusters',
                component: ClusterComponent,
            },
            {
                path: 'clusters/:name',
                component: ClusterDetailComponent,
                resolve: {cluster: ClusterRoutingResolverService},
                children: [
                    {path: '', redirectTo: 'overview', pathMatch: 'full'},
                    {path: 'overview', component: OverviewComponent},
                    {path: 'nodes', component: NodeComponent},
                    {path: 'namespaces', component: NamespaceComponent},
                    {
                        path: 'storages',
                        component: StorageComponent,
                        children: [
                            {path: '', redirectTo: 'pv', pathMatch: 'full'},
                            {path: 'pv', component: PersistentVolumeComponent},
                            {path: 'pvc', component: PersistentVolumeClaimComponent},
                        ],
                    },
                    {
                        path: 'workloads',
                        component: WorkloadComponent,
                        children: [
                            {path: '', redirectTo: 'deployment', pathMatch: 'full'},
                            {path: 'deployment', component: DeploymentComponent},
                            {path: 'statefulset', component: StatefulSetComponent},
                            {path: 'daemonset', component: DaemonSetComponent},
                            {path: 'job', component: JobComponent},
                            {path: 'cornjob', component: CornJobComponent},
                        ],
                    },
                    {
                        path: 'service',
                        component: ServiceComponent,
                    },
                    {
                        path: 'ingress',
                        component: IngressComponent,
                    },
                    {
                        path: 'config',
                        component: ConfigComponent,
                        children: [
                            {path: '', redirectTo: 'cm', pathMatch: 'full'},
                            {path: 'cm', component: ConfigMapComponent},
                            {path: 'secret', component: SecretComponent},
                        ],
                    },
                    {path: 'logging', component: LoggingComponent},
                    {path: 'monitor', component: MonitorComponent},
                    {path: 'tasks', component: TaskComponent},
                ],
            },
            {
                path: 'hosts',
                component: HostComponent,
            },
            {
                path: 'setting',
                component: SettingComponent,
                children: [
                    {path: 'credential', component: CredentialComponent}
                ]
            },
            {
                path: 'users',
                component: UserComponent,
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
