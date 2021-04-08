import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerCreateGlusterfsComponent } from './storage-provisioner-create-glusterfs.component';

describe('StorageProvisionerCreateGlusterfsComponent', () => {
  let component: StorageProvisionerCreateGlusterfsComponent;
  let fixture: ComponentFixture<StorageProvisionerCreateGlusterfsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageProvisionerCreateGlusterfsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerCreateGlusterfsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
