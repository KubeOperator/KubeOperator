import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LicenseImportComponent } from './license-import.component';

describe('LicenseImportComponent', () => {
  let component: LicenseImportComponent;
  let fixture: ComponentFixture<LicenseImportComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LicenseImportComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LicenseImportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
