import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';
import {KubernetesService} from '../../kubernetes.service';
import {V1Deployment, V1Namespace, V1Node, V1Pod} from '@kubernetes/client-node';
import {OverViewData} from '../../cluster'

@Component({
    selector: 'app-overview',
    templateUrl: './overview.component.html',
    styleUrls: ['./webkubectl/overview.component.css']
})
export class OverviewComponent implements OnInit {

    constructor(private route: ActivatedRoute,
                private kubernetesService: KubernetesService) {
    }

    currentCluster: Cluster;
    nodes: V1Node[] = [];
    overViewData: OverViewData;
    cpuTotal = 0;
    cpuUsage = 0;
    cpuUsagePercent = 0.0;
    memTotal = 0;
    memUsage = 0;
    memUsagePercent = 0.0;
    podLimit = 0;
    podUsagePercent = 0.0;


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.loadDatas();
        });
    }

    loadDatas() {
        let search = {
            kind: "overviewdatas",
            cluster: this.currentCluster.name,
            continue: "",
            limit: 0,
            namespace: "",
            name: "",
        }
        this.kubernetesService.listResource(search).subscribe(data => {
            this.overViewData = data;
            this.listNodes()
        });
    }

    listNodes() {
        let search = {
            kind: "nodelist",
            cluster: this.currentCluster.name,
            continue: "",
            limit: 0,
            namespace: "",
            name: "",
        }
        this.kubernetesService.listResource(search).subscribe(data => {
            this.nodes = data.items;
            for (const node of this.nodes) {
                this.cpuTotal = this.cpuTotal + Number(node.status.capacity.cpu);
                const mem = node.status.capacity.memory.replace('Ki', '');
                this.memTotal = this.memTotal + Number(mem);
                this.podLimit = this.podLimit + Number(node.status.capacity.pods);
            }
            this.podUsagePercent = (this.overViewData.pods / this.podLimit) * 100;
            this.listNodesUsage();
        });
    }

    listNodesUsage() {
        this.kubernetesService.getMetrics(this.currentCluster.name).subscribe(data => {
            const metrics = data.items;
            for (const me of metrics) {
                const c = me.usage.cpu.replace('n', '');
                this.cpuUsage = this.cpuUsage + Number(c);
                const m = me.usage.memory.replace('Ki', '');
                this.memUsage = this.memUsage + Number(m);
            }
            this.cpuUsage = this.cpuUsage / (1000 * 1000 * 1000);
            this.memUsagePercent = (this.memUsage / this.memTotal) * 100;
            this.cpuUsagePercent = (this.cpuUsage / this.cpuTotal) * 100;
        });
    }
}

