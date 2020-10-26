import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageProvisionerCreateOceanStorComponent } from './storage-provisioner-create-ocean-stor.component';

describe('StorageProvisionerCreateOceanStorComponent', () => {
  let component: StorageProvisionerCreateOceanStorComponent;
  let fixture: ComponentFixture<StorageProvisionerCreateOceanStorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ StorageProvisionerCreateOceanStorComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageProvisionerCreateOceanStorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
