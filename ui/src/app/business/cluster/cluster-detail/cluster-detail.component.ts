import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {Cluster} from '../cluster';
import {ToolsService} from './tools/tools.service';
import {ClusterTool} from './tools/tools';
import {BusinessLicenseService} from '../../../shared/service/business-license.service';

@Component({
    selector: 'app-cluster-detail',
    templateUrl: './cluster-detail.component.html',
    styleUrls: ['./cluster-detail.component.css']
})
export class ClusterDetailComponent implements OnInit {

    constructor(private router: Router,
                private route: ActivatedRoute,
                private toolsService: ToolsService,
                private businessLicenseService: BusinessLicenseService) {
    }

    currentCluster: Cluster;
    tools: ClusterTool[] = [];
    ready = false;
    hasLicense = false;

    ngOnInit(): void {
        this.route.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.toolsService.list(this.currentCluster.name).subscribe(d => {
                if (d) {
                    this.tools = d;
                }
                this.ready = true;
            });
        });
        this.hasLicense = this.businessLicenseService.licenseValid;
    }

    showApp(toolName: string) {
        for (const tool of this.tools) {
            if (tool.name === toolName && tool.status === 'Running') {
                return true;
            }
        }
        return false;
    }

    backToCluster() {
        this.router.navigate(['projects/' + this.currentCluster.projectName]).then();
    }

}
