import {Component, Input, OnInit} from '@angular/core';
import {Cluster} from "../../../cluster";
import {WebkubectlService} from "./webkubectl.service";
import {WebkubectlToken} from "./webkubectl";
import {DomSanitizer} from "@angular/platform-browser";

@Component({
    selector: 'app-webkubectl',
    templateUrl: './webkubectl.component.html',
    styleUrls: ['./webkubectl.component.css']
})
export class WebkubectlComponent implements OnInit {

    constructor(private service: WebkubectlService, private sanitizer: DomSanitizer) {
    }

    @Input() currentCluster: Cluster;
    url: any;
    loading = false;
    opened = false;

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
        this.loading = true;
        this.service.getToken(this.currentCluster.name).subscribe(data => {
            this.url = this.sanitizer.bypassSecurityTrustResourceUrl('/webkubectl/terminal/?token=' + data.token);
            this.loading = false;
        });
    }

    newWindow() {
        this.opened = true;
        this.loading = true;
        this.service.getToken(this.currentCluster.name).subscribe(data => {
            this.url = `/webkubectl/terminal/?token=${data.token}`
            this.loading = false;
            window.open(this.url,'_blank','weight=300,height=200,alwaysRaised=yes,depended=yes')
        });
    }


}
