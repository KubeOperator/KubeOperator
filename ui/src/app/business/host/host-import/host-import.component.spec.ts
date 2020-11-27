import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HostImportComponent } from './host-import.component';

describe('HostImportComponent', () => {
  let component: HostImportComponent;
  let fixture: ComponentFixture<HostImportComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ HostImportComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(HostImportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
