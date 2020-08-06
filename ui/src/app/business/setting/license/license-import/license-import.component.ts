import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {NgForm} from "@angular/forms";
import {HttpClient} from "@angular/common/http";

@Component({
    selector: 'app-license-import',
    templateUrl: './license-import.component.html',
    styleUrls: ['./license-import.component.css']
})
export class LicenseImportComponent implements OnInit {

    constructor(private http: HttpClient) {
    }

    opened = false;
    file;
    @ViewChild('itemForm') itemForm: NgForm;
    @Output() imported = new EventEmitter();

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
    }

    upload(e): void {
        this.file = e.target.files[0];
    }

    onSubmit() {
        const formData = new FormData();
        formData.append('file', this.file);
        this.http.post('/api/v1/license', formData).subscribe(data => {
            this.opened = false;
            this.imported.emit();
        });
    }

}
