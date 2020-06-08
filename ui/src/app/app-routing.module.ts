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

const routes: Routes = [
    {path: 'login', component: LoginComponent},
    {
        path: '',
        component: LayoutComponent,
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
                    }
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
