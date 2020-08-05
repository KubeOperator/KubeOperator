import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerCreateVsphereComponent } from './storage-provisioner-create-vsphere.component';

describe('StorageProvisionerCreateVsphereComponent', () => {
  let component: StorageProvisionerCreateVsphereComponent;
  let fixture: ComponentFixture<StorageProvisionerCreateVsphereComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerCreateVsphereComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerCreateVsphereComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
