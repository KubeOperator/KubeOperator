import {NgModule} from '@angular/core';
import {ClusterComponent} from './cluster.component';
import {CoreModule} from '../../core/core.module';
import {ClusterListComponent} from './cluster-list/cluster-list.component';
import {ClusterCreateComponent} from './cluster-create/cluster-create.component';
import {ClusterDeleteComponent} from './cluster-delete/cluster-delete.component';
import {ClusterDetailComponent} from './cluster-detail/cluster-detail.component';
import {OverviewComponent} from './cluster-detail/overview/overview.component';
import {RouterModule} from '@angular/router';
import {ClusterConditionComponent} from './cluster-condition/cluster-condition.component';
import {NodeComponent} from './cluster-detail/node/node.component';
import {ServiceComponent} from './cluster-detail/service/service.component';
import {IngressComponent} from './cluster-detail/ingress/ingress.component';
import {NamespaceComponent} from './cluster-detail/namespace/namespace.component';
import {NamespaceListComponent} from './cluster-detail/namespace/namespace-list/namespace-list.component';
import {SharedModule} from '../../shared/shared.module';
import {StorageComponent} from './cluster-detail/storage/storage.component';
import {PersistentVolumeComponent} from './cluster-detail/storage/persistent-volume/persistent-volume.component';
import {PersistentVolumeClaimComponent} from './cluster-detail/storage/persistent-volume-claim/persistent-volume-claim.component';
import {StorageClassComponent} from './cluster-detail/storage/storage-class/storage-class.component';
import {PersistentVolumeListComponent} from './cluster-detail/storage/persistent-volume/persistent-volume-list/persistent-volume-list.component';
import {NodeListComponent} from './cluster-detail/node/node-list/node-list.component';
import {NodeDetailComponent} from './cluster-detail/node/node-detail/node-detail.component';
import {ConfigComponent} from './cluster-detail/config/config.component';
import {PersistentVolumeClaimListComponent} from './cluster-detail/storage/persistent-volume-claim/persistent-volume-claim-list/persistent-volume-claim-list.component';
import {StorageClassListComponent} from './cluster-detail/storage/storage-class/storage-class-list/storage-class-list.component';
import {WorkloadComponent} from './cluster-detail/workload/workload.component';
import {DeploymentComponent} from './cluster-detail/workload/deployment/deployment.component';
import {StatefulSetComponent} from './cluster-detail/workload/stateful-set/stateful-set.component';
import {DaemonSetComponent} from './cluster-detail/workload/daemon-set/daemon-set.component';
import {JobComponent} from './cluster-detail/workload/job/job.component';
import {CornJobComponent} from './cluster-detail/workload/corn-job/corn-job.component';
import {DeploymentListComponent} from './cluster-detail/workload/deployment/deployment-list/deployment-list.component';
import {StatefulSetListComponent} from './cluster-detail/workload/stateful-set/stateful-set-list/stateful-set-list.component';
import {JobListComponent} from './cluster-detail/workload/job/job-list/job-list.component';
import {CornJobListComponent} from './cluster-detail/workload/corn-job/corn-job-list/corn-job-list.component';
import {DaemonSetListComponent} from './cluster-detail/workload/daemon-set/daemon-set-list/daemon-set-list.component';
import {ServiceListComponent} from './cluster-detail/service/service-list/service-list.component';
import {IngressListComponent} from './cluster-detail/ingress/ingress-list/ingress-list.component';
import {ConfigMapComponent} from './cluster-detail/config/config-map/config-map.component';
import {SecretComponent} from './cluster-detail/config/secret/secret.component';
import {ConfigMapListComponent} from './cluster-detail/config/config-map/config-map-list/config-map-list.component';
import {SecretListComponent} from './cluster-detail/config/secret/secret-list/secret-list.component';
import {LoggingComponent} from './cluster-detail/logging/logging.component';
import {LoggingQueryComponent} from './cluster-detail/logging/logging-query/logging-query.component';
import {TaskComponent} from './cluster-detail/task/task.component';
import {TaskListComponent} from './cluster-detail/task/task-list/task-list.component';
import {IngressInstallComponent} from './cluster-detail/ingress/ingress-install/ingress-install.component';


@NgModule({
    declarations: [ClusterComponent, ClusterListComponent, ClusterCreateComponent, ClusterDeleteComponent, ClusterDetailComponent,
        OverviewComponent, ClusterConditionComponent, NodeComponent, NamespaceComponent, NamespaceListComponent,
        StorageComponent, PersistentVolumeComponent, PersistentVolumeClaimComponent, StorageClassComponent, PersistentVolumeListComponent,
        NodeListComponent, NodeDetailComponent, ConfigComponent, PersistentVolumeClaimListComponent,
        StorageClassListComponent, WorkloadComponent, DeploymentComponent, StatefulSetComponent, DaemonSetComponent, JobComponent,
        CornJobComponent, DeploymentListComponent, StatefulSetListComponent, JobListComponent, CornJobListComponent,
        DaemonSetListComponent, ServiceComponent, IngressComponent, ServiceListComponent, IngressListComponent, ConfigMapComponent,
        SecretComponent, ConfigMapListComponent, SecretListComponent, LoggingComponent, LoggingQueryComponent, TaskComponent,
        TaskListComponent, IngressInstallComponent],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
    ]
})
export class ClusterModule {
}
